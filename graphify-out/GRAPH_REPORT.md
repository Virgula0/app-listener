# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 499 files · ~1,723,172 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6651 nodes · 28221 edges · 259 communities (201 shown, 58 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2662 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f5a8dbdc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_os
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- policyEditor
- protectedPage
- fileedit_symlink_test.go
- plan.go
- safeio.go
- IntegrationSuite
- guard_trust.bpf.c
- time.Time
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- pinstate.go
- guard.bpf.c
- TrustGuard
- Finder
- codeedit.go
- shellrc.go
- BinaryEntry
- TestMonitorUseCaseLifecycle
- UI
- github.com/charmbracelet/bubbletea.Cmd
- picker
- GenerateConf
- IntegrationSuite
- Message
- controlServer
- EventType
- NetEventType
- IntegrationSuite
- absPath
- monitorUnitTest
- gui.go
- fileedit_model_test.go
- User
- stdlib
- lockAndDeprovision
- fadeScroll
- migrate_test.go
- fyne.io/fyne/v2/driver/desktop.Cursor
- Client
- daemonconfig.go
- fyne.io/fyne/v2/widget.Label
- Plan
- FileEvent
- mutedText
- engine
- os.File
- New
- Target
- inode_dev
- installPage
- runDaemon
- sync.Mutex
- gui/session.go
- IntegrationSuite
- dialLiveSession
- update_test.go
- Network
- pickfs.go
- Guard
- stdio
- bpfstats/main.go
- IntegrationSuite
- RefreshCatalog
- NetGuardEvent
- watchSet
- process_vm_readv.c
- .readEvent
- editView
- planUpdaters
- update/update.go
- net_tester/main.go
- uninstall_system.go
- recordFlow
- go_pkg_strings
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runLockdown
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- github.com/charmbracelet/bubbles/textarea.Model
- Well-known bypass classes
- install/install.go
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- IntegrationSuite
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- installcert.go
- mountTable
- procstats.go
- BinaryStat
- IntegrationSuite
- io.Writer
- openFileVaultWithMasterKey
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- RefreshSection
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- newTestTarget
- updatePage
- trace-app-libs.sh
- IntegrationSuite
- root.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- runNetworkMonitor
- serveWebSocket
- selection_test.go
- newTestVetter
- Ledger
- trustManager
- recoverFileInPlace
- fscrypt.go
- runMonitor
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- NetworkMonitor
- NewDumpHook
- hashCache
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- vaultfs.go
- .populateWatchInodes
- feed.go
- shell
- dialog_test.go
- audit.go
- SanitizeText
- WithBottomBar
- DumpHook
- Limitations
- Troubleshooting
- IntegrationSuite
- appTheme
- netinject.c
- copyRegularAt
- Binary
- Compatibility
- How it works
- github.com/spf13/cobra.Command
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- changelog.go
- IntegrationSuite
- edit-protected
- monitor.md
- kernel-prof.sh
- deploy
- pathCache
- inspector_probe.c
- DaemonEvent
- AddServeFlags
- github.com/stretchr/testify/suite.Suite
- finishBinary
- Desktop GUI
- IntegrationSuite
- ListUsers
- run
- .eventTypeAllowed
- bleedLayout
- Classify
- IntegrationSuite
- .awaitDaemonUpOrDead
- LibraryClosure
- KernelDev
- CleanupStalePins
- NewEventFanout
- TestCleanRel
- Existing
- prctl_gate
- LoaderRefusal
- tempGrantSetup
- supersede_probe.c
- Prereq
- downloadFile
- Vault
- encryptDirectories
- Execute
- taint_victim.c
- .runEditorHarness
- serve.go
- configedit_test.go
- RawTarget
- GuardInodeKey
- Daemon configuration
- serveHandler
- WithSuspectsFrom
- loadDaemonConfig
- parseResize

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `shell` - 91 edges
4. `Load()` - 89 edges
5. `absPath()` - 88 edges
6. `IntegrationSuite` - 83 edges
7. `fileEditModel` - 77 edges
8. `writeConfig()` - 69 edges
9. `IntegrationSuite` - 68 edges
10. `GuardInodeKey` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `certifyInstall()` --calls--> `InstallCertPath()`  [EXTRACTED]
  cmd/functions/install/install.go → internal/protected/installcert.go

## Import Cycles
- None detected.

## Communities (259 total, 58 thin omitted)

### Community 0 - "go_pkg_os"
Cohesion: 0.16
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (114): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped(), TestForwardFlagSurface(), TestSubcommandsRegistered() (+106 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.19
Nodes (12): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (74): TestBuildTrustedSetIncludesPendingBinaries(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+66 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (28): addErrKind, bpfGuardEvent, ReplacementCheck, BinariesSummary(), classifyAddErr(), decodeGuardEvent(), fsGateLabel(), guardModeKey() (+20 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "GuardEvent"
Cohesion: 0.05
Nodes (17): holdHeadless(), HeadlessLine(), snapEngines, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+9 more)

### Community 11 - "policyEditor"
Cohesion: 0.11
Nodes (13): handRadioGroup, netGuardTab, policyEditor, policyHooks, policyTab, newRadioGroup(), newPolicyTab(), verbFor() (+5 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (22): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordForm, recordReview (+14 more)

### Community 13 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 14 - "plan.go"
Cohesion: 0.28
Nodes (10): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), recordTick, linkedFile, parsePasswd(), shell(), TestParsePasswd() (+2 more)

### Community 15 - "safeio.go"
Cohesion: 0.08
Nodes (36): checkAndApply(), refreshCatalogConfig(), dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent() (+28 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "time.Time"
Cohesion: 0.07
Nodes (12): LockedError, sweepID, sweepMemo, logGrid, logLine, runView, logRow(), logStyle() (+4 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.10
Nodes (61): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+53 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (17): trustHook, trustRows, ensureBtrfsLayout(), deleteInoKeys(), deleteResKeys(), setSupersedeTrust(), loadPinnedExeMaps(), stripPinnedAllow() (+9 more)

### Community 27 - "Finder"
Cohesion: 0.12
Nodes (21): findButton, span, column(), reposition(), findAll(), Finder, lowerRunes(), NewFinder() (+13 more)

### Community 28 - "codeedit.go"
Cohesion: 0.09
Nodes (41): class, classOf(), Edit(), leadingIndent(), moveTo(), Navigate(), New(), onlySpaces() (+33 more)

### Community 29 - "shellrc.go"
Cohesion: 0.16
Nodes (35): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+27 more)

### Community 30 - "BinaryEntry"
Cohesion: 0.08
Nodes (21): keptEntries, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), WithEagerPopulate(), diffBinaries(), appendByte(), guardUnitTest (+13 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+4 more)

### Community 32 - "UI"
Cohesion: 0.08
Nodes (22): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, EditText() (+14 more)

### Community 33 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.05
Nodes (26): recordTickCmd(), recordModel, diffModel, editorModel, formatType(), listenForEvents(), syncViewport(), tickStats() (+18 more)

### Community 34 - "picker"
Cohesion: 0.12
Nodes (7): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry(), newPickRow(), invokingHome()

### Community 35 - "GenerateConf"
Cohesion: 0.11
Nodes (24): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig() (+16 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (11): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), File, backend, backendOp, backendState, Message (+3 more)

### Community 38 - "controlServer"
Cohesion: 0.11
Nodes (20): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+12 more)

### Community 39 - "EventType"
Cohesion: 0.05
Nodes (30): fakeDaemon, exeRow, tempGrant, tempMaskOp, tempRow, EventType, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+22 more)

### Community 40 - "NetEventType"
Cohesion: 0.11
Nodes (21): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), handCheckGroup, NetEventType (+13 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (8): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 44 - "gui.go"
Cohesion: 0.17
Nodes (16): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode (+8 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "User"
Cohesion: 0.08
Nodes (28): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchersFor(), globBuilder (+20 more)

### Community 48 - "lockAndDeprovision"
Cohesion: 0.12
Nodes (10): entryOf(), deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy() (+2 more)

### Community 49 - "fadeScroll"
Cohesion: 0.05
Nodes (27): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "migrate_test.go"
Cohesion: 0.10
Nodes (20): stripEntryFlags(), stripImmutableFlags(), TestDecryptFileRefusesNonEncrypted(), TestDecryptMissingPath(), TestDecryptRefusesExistingTempDir(), TestDecryptRefusesPlainDirectory(), TestEncryptFileRefusesExistingBackup(), TestEncryptFileRefusesHardlink() (+12 more)

### Community 51 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.15
Nodes (5): buttonCursor, handSelect, handCursor(), newButtonCursor(), TestButtonCursorFollowsState()

### Community 52 - "Client"
Cohesion: 0.08
Nodes (16): Client, Reader, Writer, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend() (+8 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (48): systemPatterns(), systemRule(), grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource() (+40 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (23): binaryRow, clickPanel, counterStrip, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+15 more)

### Community 55 - "Plan"
Cohesion: 0.11
Nodes (18): BinaryRule, Plan, materialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan() (+10 more)

### Community 56 - "FileEvent"
Cohesion: 0.12
Nodes (5): snapMonitor, FileEvent, fakeSource, MonitorRepository, fakeMonitorRepo

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (48): binaryList, checkGrid, filesPage, handButton, handCheck, homePage, multiAnswer, panel (+40 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (42): markReady(), notifySystemdReady(), engine, inodeBatchFn, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank() (+34 more)

### Community 59 - "os.File"
Cohesion: 0.05
Nodes (63): setConfinedHomes(), forwardCheck(), VettedInode, Multicall, MulticallKind, verifiedBinary(), ExeKey(), IsMulticallRefusal() (+55 more)

### Community 60 - "New"
Cohesion: 0.23
Nodes (17): collectFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails() (+9 more)

### Community 61 - "Target"
Cohesion: 0.10
Nodes (31): Target, appletLine(), classifyGroup(), classifyLibs(), innermostResource(), libCovered(), newBinary(), verifyExe() (+23 more)

### Community 62 - "inode_dev"
Cohesion: 0.11
Nodes (29): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+21 more)

### Community 63 - "installPage"
Cohesion: 0.07
Nodes (19): actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, Options (+11 more)

### Community 64 - "runDaemon"
Cohesion: 0.06
Nodes (25): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newControlManager(), startControlManager() (+17 more)

### Community 65 - "sync.Mutex"
Cohesion: 0.08
Nodes (18): guardPill(), newActivity(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget() (+10 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "dialLiveSession"
Cohesion: 0.10
Nodes (23): reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), editAgain(), fetchConfig(), putConfig() (+15 more)

### Community 69 - "update_test.go"
Cohesion: 0.09
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 70 - "Network"
Cohesion: 0.12
Nodes (5): netGuardPill(), clonePolicy(), NetPolicy, Network, NetExtra

### Community 71 - "pickfs.go"
Cohesion: 0.24
Nodes (13): listOpts, pickEntry, treeSearch, TestSearchTreeFollowsNoSymlink(), plural(), entryIsDir(), filterEntries(), listDir() (+5 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (22): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, SectionScan, openRoot() (+14 more)

### Community 73 - "stdio"
Cohesion: 0.06
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.12
Nodes (25): confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad() (+17 more)

### Community 75 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 76 - "RefreshCatalog"
Cohesion: 0.14
Nodes (17): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+9 more)

### Community 77 - "NetGuardEvent"
Cohesion: 0.09
Nodes (10): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NetworkGuardUseCase, NewNetworkGuardUseCase(), TestNetworkGuardUseCaseStartError(), netGuardEventLine, fakeNet (+2 more)

### Community 78 - "watchSet"
Cohesion: 0.15
Nodes (13): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - ".readEvent"
Cohesion: 0.15
Nodes (9): BpfEvent, eventUnitTest, NetBpfEvent, Cstr(), DecodeBpfEvent(), FormatAddr(), Ntohs(), reasonText() (+1 more)

### Community 81 - "editView"
Cohesion: 0.08
Nodes (17): codeEntry, editView, Entry, selectable(), handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm() (+9 more)

### Community 82 - "planUpdaters"
Cohesion: 0.08
Nodes (26): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+18 more)

### Community 83 - "update/update.go"
Cohesion: 0.23
Nodes (18): newReleaseWindow(), releasesSince(), channelOf(), compareStableVersions(), filterChannel(), findUpdate(), hasAsset(), Latest() (+10 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "uninstall_system.go"
Cohesion: 0.08
Nodes (37): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+29 more)

### Community 86 - "recordFlow"
Cohesion: 0.17
Nodes (14): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+6 more)

### Community 87 - "go_pkg_strings"
Cohesion: 0.08
Nodes (9): Progress, EnsurePlaceholders(), JournalPath(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), DeleteBackup(), TestDeleteBackupOnlyItsOwn(), i3Result (+1 more)

### Community 88 - "liveSession"
Cohesion: 0.24
Nodes (9): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward() (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runLockdown"
Cohesion: 0.15
Nodes (15): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs() (+7 more)

### Community 92 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 101 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (19): lineStyles, pos, viewportState, dedent(), lines(), baseStyle(), gutterLabel(), lineLen() (+11 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "install/install.go"
Cohesion: 0.06
Nodes (67): offerDaemonStart(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), deploy() (+59 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "installcert.go"
Cohesion: 0.11
Nodes (31): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+23 more)

### Community 111 - "mountTable"
Cohesion: 0.16
Nodes (10): MountInfo, mountinfoUnitTest, TestNoteRootMountFollowsMountID(), MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo() (+2 more)

### Community 112 - "procstats.go"
Cohesion: 0.07
Nodes (30): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+22 more)

### Community 113 - "BinaryStat"
Cohesion: 0.11
Nodes (22): BinaryStat, placedWalk, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), confineBelow(), inUserHome() (+14 more)

### Community 115 - "io.Writer"
Cohesion: 0.31
Nodes (8): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, candidate

### Community 116 - "openFileVaultWithMasterKey"
Cohesion: 0.18
Nodes (14): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "RefreshSection"
Cohesion: 0.07
Nodes (42): User, LibraryBlock, RefreshOptions, AppendSectionWhitelist(), ConfSafePath(), findSectionEnd(), findSectionStart(), firstContentLine() (+34 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "newTestTarget"
Cohesion: 0.13
Nodes (20): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+12 more)

### Community 131 - "updatePage"
Cohesion: 0.16
Nodes (6): latestRelease, statusTile, updatePage, newUpdatePage(), boolIndex(), isPreRelease()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "root.go"
Cohesion: 0.12
Nodes (8): releaseBuild(), Run(), stageItem, stageState, textAnswer, encodeTable(), TestEncodeTable(), rgb()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "runNetworkMonitor"
Cohesion: 0.18
Nodes (8): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NewNetModel(), NetworkMonitorUseCase

### Community 143 - "serveWebSocket"
Cohesion: 0.21
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 144 - "selection_test.go"
Cohesion: 0.25
Nodes (10): clipboardText(), setClipboard(), cellOf(), repeat(), selKey(), TestSelectionEditing(), TestSelectionLargeDelete(), TestSelectionMouseDrag() (+2 more)

### Community 145 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "trustManager"
Cohesion: 0.14
Nodes (17): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+9 more)

### Community 148 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 149 - "fscrypt.go"
Cohesion: 0.05
Nodes (62): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+54 more)

### Community 150 - "runMonitor"
Cohesion: 0.18
Nodes (10): MakeDisplayPaths(), ParseEventsFlag(), validateForwardEvents(), runHeadless(), runMonitor(), runTUI(), NewModel(), MonitorUseCase (+2 more)

### Community 152 - "inodeChain"
Cohesion: 0.26
Nodes (13): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+5 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.05
Nodes (53): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration() (+45 more)

### Community 155 - "NetworkMonitor"
Cohesion: 0.12
Nodes (7): snapNetMonitor, NetEvent, NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 156 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 157 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 163 - "vaultfs.go"
Cohesion: 0.07
Nodes (36): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+28 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (37): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, networkPage (+29 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (50): daemonPage, daemonState, fakeLink, pageDef, pageID, reel, resizeWatch, runHost (+42 more)

### Community 167 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 168 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 169 - "SanitizeText"
Cohesion: 0.11
Nodes (13): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), startTrustGuard(), trustStartupError(), binaryVetter, vetDecision (+5 more)

### Community 170 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 171 - "DumpHook"
Cohesion: 0.24
Nodes (5): logHook, formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 179 - "Binary"
Cohesion: 0.17
Nodes (14): formatElapsed(), lines(), orNone(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording() (+6 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.20
Nodes (15): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials(), TestSetupDumpLogRefusesSymlink() (+7 more)

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "changelog.go"
Cohesion: 0.22
Nodes (12): changelogText(), cumulativeChangelog(), newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize() (+4 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "deploy"
Cohesion: 0.14
Nodes (24): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy() (+16 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 202 - "DaemonEvent"
Cohesion: 0.11
Nodes (31): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+23 more)

### Community 203 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 204 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.14
Nodes (11): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), mustParse(), TestValidateTargets() (+3 more)

### Community 205 - "finishBinary"
Cohesion: 0.28
Nodes (8): beyondEvents(), finishBinary(), isShellHelper(), noteWhitelisted(), TestFileReadersStartDeselected(), TestNoteWhitelisted(), TestShellHelpersStartDeselected(), exeVerdict

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "ListUsers"
Cohesion: 0.07
Nodes (50): askBunTmpdirUsers(), gatherPerUserSetup(), addedResources(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths(), pathCovered() (+42 more)

### Community 210 - "run"
Cohesion: 0.24
Nodes (11): deployGUI(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), RunWith(), syncGUI() (+3 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "Classify"
Cohesion: 0.32
Nodes (6): Classify(), fixtureExe(), privateTempRoots(), TestClassifyLibs(), TestClassifyRefusesWorldWritableLocations(), TestClassifyVerifiesTheInodeThatRan()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 217 - "KernelDev"
Cohesion: 0.11
Nodes (23): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+15 more)

### Community 218 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 219 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "LoaderRefusal"
Cohesion: 0.67
Nodes (3): loaderName(), LoaderRefusal(), soNames()

### Community 224 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "Prereq"
Cohesion: 0.29
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 228 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 230 - "encryptDirectories"
Cohesion: 0.38
Nodes (4): encryptDirectories(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 232 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 237 - "serve.go"
Cohesion: 0.23
Nodes (9): isInteractiveTerminal(), newServeListener(), relayServeSignals(), runServe(), Serve(), serveHTTP(), waitProgram(), resizeMessage (+1 more)

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 241 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 242 - "GuardInodeKey"
Cohesion: 0.08
Nodes (23): refusedReplacements, supersedeUnnamed(), appletTag(), buildOf(), GuardInodeKey, mcNameKey(), MulticallNameTags(), registerBuild() (+15 more)

### Community 244 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 245 - "serveHandler"
Cohesion: 0.27
Nodes (9): basicAuth(), newServeServer(), securityHeaders(), serveHandler(), TestBasicAuth(), TestSecurityHeadersAndExactRoot(), TestServeHandlerProtectsPageAndWebSocket(), TestServeWebSocketSnapshotsResizeAndSingleViewer() (+1 more)

### Community 247 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 248 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 252 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 658 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `Message`, `installcert.go`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025777777777777778 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._
- **Why does `fakeVault` connect `Prereq` to `plan.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07403508771929825 - nodes in this community are weakly interconnected._