# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 488 files · ~1,716,925 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6513 nodes · 27935 edges · 260 communities (203 shown, 57 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2602 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4a8a7d85`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
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
- highlight_test.go
- procstats.go
- editView
- IntegrationSuite
- guard_trust.bpf.c
- runView
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
- rcUser
- catalogRefresher
- NetGuardEvent
- Active
- runNetworkGuard
- picker
- dialLiveSession
- time.Duration
- Message
- controlServer
- TempBinary
- IntegrationSuite
- absPath
- monitorUnitTest
- RunWith
- WithBottomBar
- User
- stdio
- CopyTreeWithProgress
- fadeScroll
- newTestTarget
- cursor.go
- Client
- watchGroup
- tabs.go
- .apply
- runMonitor
- mutedText
- GuardInodeKey
- os.File
- daemon_multicall_test.go
- Recorder
- exe_supersede.h
- Network
- Config
- FileTarget
- gui/session.go
- IntegrationSuite
- os.FileMode
- update_test.go
- unistd
- NetworkMonitor
- Guard
- string
- SampleContent
- RefreshSection
- journal
- fakeGuardRepo
- process_vm_readv.c
- compileGlobNames
- update/update.go
- planUpdaters
- changelog.go
- net_tester/main.go
- runGuard
- selection.go
- recordFlow
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- fileedit_model_test.go
- Monitor
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- pulse_test.go
- installPage
- github.com/charmbracelet/bubbles/textarea.Model
- Well-known bypass classes
- UI
- jit_provenance.c
- inspector_probe.c
- recordops.c
- KernelDev
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- bpfstats/main.go
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- ConfinedEntry
- IntegrationSuite
- buildTrustedSet
- .FilterExistingWhitelist
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- .launch
- supersede_probe.c
- Candidate
- trace-app-libs.sh
- IntegrationSuite
- ListUsers
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- MountVouch
- BinaryEntry
- DeployInstalledBinary
- newTestVetter
- Ledger
- fileedit_symlink_test.go
- New
- SelectOrphans
- netGuardModel
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- headless.go
- NewDumpHook
- runBackend
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- .step
- Vault
- EventType
- sync.Mutex
- shell
- ValidateTargets
- SanitizeText
- taint_victim.c
- netModel
- .populateWatchInodes
- Limitations
- Troubleshooting
- IntegrationSuite
- hashCache
- model
- copyRegularAt
- find_test.go
- finishBinary
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
- secureResources
- IntegrationSuite
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- fetchReleases
- pathCache
- .runEditorHarness
- DaemonEvent
- runNetworkMonitor
- daemon: permanent protection with encryption at rest
- Desktop GUI
- IntegrationSuite
- applyDiffAdditions
- bpfLsmListed
- .eventTypeAllowed
- diffModel
- Prereq
- selection_test.go
- LibraryClosure
- Vault
- TestCleanRel
- Daemon configuration
- unlockUnderGuard
- time.Time
- main_test.go
- tempGrantSetup
- changelog_test.go
- IntegrationSuite
- fakeVault
- serve.go
- configedit_test.go
- writeWithin
- openBunTmp
- Execute
- fakeDaemon
- ReadInstallCert
- record
- lockAndDeprovision
- github.com/google/fscrypt/actions.Context
- .exeKeys
- vouchedDevs
- testRecording
- installHooks
- parseKernelRelease
- copyXattrs
- golang.org/x/sys/unix.Timespec
- MulticallError
- isNotEncryptedErr

## God Nodes (most connected - your core abstractions)
1. `Config` - 109 edges
2. `Guard` - 91 edges
3. `Load()` - 89 edges
4. `shell` - 87 edges
5. `absPath()` - 84 edges
6. `IntegrationSuite` - 82 edges
7. `GuardInodeKey` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (260 total, 57 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (19): tempJournalRow, Progress, User, invalidPassword, recordTick, launchRule, trustEvent, textAnswer (+11 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (148): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), loadDaemonConfig(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown() (+140 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.16
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

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
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (47): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+39 more)

### Community 10 - "GuardEvent"
Cohesion: 0.11
Nodes (15): bpfGuardEvent, GuardEvent, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary(), decodeGuardEvent(), fsGateLabel() (+7 more)

### Community 11 - "policyEditor"
Cohesion: 0.12
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (26): backendLink, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage, recordForm (+18 more)

### Community 13 - "highlight_test.go"
Cohesion: 0.16
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 14 - "procstats.go"
Cohesion: 0.11
Nodes (20): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+12 more)

### Community 15 - "editView"
Cohesion: 0.11
Nodes (10): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+2 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (15): logGrid, logHook, logLine, runHost, runView, stageItem, stageState, logRow() (+7 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), chownUser, dirEntry, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseReloadAddedAlreadyProvisionedNotUnlocked() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.08
Nodes (11): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, setSupersedeTrust(), liveTrust() (+3 more)

### Community 27 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "codeedit.go"
Cohesion: 0.18
Nodes (27): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+19 more)

### Community 29 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 30 - "catalogRefresher"
Cohesion: 0.10
Nodes (12): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer, NewEventFanout() (+4 more)

### Community 31 - "NetGuardEvent"
Cohesion: 0.08
Nodes (19): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo() (+11 more)

### Community 32 - "Active"
Cohesion: 0.20
Nodes (10): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), Active() (+2 more)

### Community 33 - "runNetworkGuard"
Cohesion: 0.12
Nodes (21): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom(), DiscoverInfraBinaries() (+13 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "dialLiveSession"
Cohesion: 0.08
Nodes (29): auditAfterEdit(), auditFindings(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+21 more)

### Community 36 - "time.Duration"
Cohesion: 0.14
Nodes (9): procStats, statsSample, IntegrationSuite, pooledContainer, formatAgo(), cpuTime(), humanBytes(), rssBytes() (+1 more)

### Community 37 - "Message"
Cohesion: 0.11
Nodes (9): checkSecret(), File, backend, backendOp, backendState, BinaryCheck, Message, PwState (+1 more)

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (27): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+19 more)

### Community 39 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "RunWith"
Cohesion: 0.12
Nodes (27): runDiffCatalog(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), parseMaintenanceFlags(), runMaintenanceMode() (+19 more)

### Community 45 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 46 - "User"
Cohesion: 0.06
Nodes (41): systemPatterns(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured() (+33 more)

### Community 48 - "CopyTreeWithProgress"
Cohesion: 0.14
Nodes (10): entryOf(), migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress() (+2 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (16): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newFadeScroll(), newTestFadeScroll(), settle() (+8 more)

### Community 50 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 51 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (17): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+9 more)

### Community 53 - "watchGroup"
Cohesion: 0.17
Nodes (24): watchGroup, applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary(), applyLibDir() (+16 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (22): appTheme, binaryRow, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pill (+14 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "runMonitor"
Cohesion: 0.11
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+10 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (45): binaryList, checkGrid, clickPanel, filesPage, homePage, panel, statusTile, newLogo() (+37 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (49): engine, refusedReplacements, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+41 more)

### Community 59 - "os.File"
Cohesion: 0.04
Nodes (64): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), setConfinedHomes(), confirmUserPlaced(), forwardCheck(), VettedInode, Multicall (+56 more)

### Community 60 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 61 - "Recorder"
Cohesion: 0.13
Nodes (18): Target, Classify(), drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe() (+10 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 63 - "Network"
Cohesion: 0.07
Nodes (17): snapNetMonitor, NetEvent, NetEventType, eventsOrAll(), netGuardPill(), newNetEventChecks(), NetEventTypes(), ParseNetEventType() (+9 more)

### Community 64 - "Config"
Cohesion: 0.06
Nodes (51): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError() (+43 more)

### Community 65 - "FileTarget"
Cohesion: 0.15
Nodes (9): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+1 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "os.FileMode"
Cohesion: 0.11
Nodes (23): dupFD(), installFile(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+15 more)

### Community 69 - "update_test.go"
Cohesion: 0.07
Nodes (22): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchGUIWithoutAsset(), TestFetchReleaseByTag(), TestFilterChannel(), testKeyPair() (+14 more)

### Community 70 - "unistd"
Cohesion: 0.11
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "NetworkMonitor"
Cohesion: 0.09
Nodes (14): NetBpfEvent, appletTag(), MulticallNameTags(), TestAppletTagStablePerName(), Cstr(), FormatAddr(), Ntohs(), reasonText() (+6 more)

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (22): binaryVerifyState, deferredBinary, Guard, Mode, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), canonicalPaths() (+14 more)

### Community 73 - "string"
Cohesion: 0.11
Nodes (9): do_read(), main(), read_stat(), swap_mm_map(), put(), wait_then_read(), dump_via_debugfs(), main() (+1 more)

### Community 74 - "SampleContent"
Cohesion: 0.10
Nodes (24): pathCovered(), TestPathCovered(), TestDaemonUnitWithMetadataOutput(), isInsidePath(), TestAskSSHAgentUsersNoGuardedSSH(), askSSHAgentUsers(), detectSSHAgentUnits(), isInstallerSSHAgentUnit() (+16 more)

### Community 76 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 77 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 79 - "process_vm_readv.c"
Cohesion: 0.27
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 81 - "update/update.go"
Cohesion: 0.13
Nodes (33): archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), downloadAndVerify() (+25 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "runGuard"
Cohesion: 0.08
Nodes (26): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck(), runLockdown(), runOneShotMode() (+18 more)

### Community 86 - "selection.go"
Cohesion: 0.15
Nodes (10): pos, reposition(), cellToRune(), clipboardText(), cursorPos(), Selection, offset(), setClipboard() (+2 more)

### Community 87 - "recordFlow"
Cohesion: 0.12
Nodes (26): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries(), pickEncryption() (+18 more)

### Community 88 - "liveSession"
Cohesion: 0.18
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (26): newFileEditModel(), RunFileEditor(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow() (+18 more)

### Community 92 - "Monitor"
Cohesion: 0.16
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "netinject.c"
Cohesion: 0.19
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (20): foreignPinnedLink(), FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs(), codeVouched(), isTraced(), judgePasses() (+12 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "pulse_test.go"
Cohesion: 0.42
Nodes (6): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Sub

### Community 100 - "installPage"
Cohesion: 0.06
Nodes (21): actionGate, binSource, installAction, installedBinary, installPage, latestRelease, uninstallCard, updatePage (+13 more)

### Community 101 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "UI"
Cohesion: 0.12
Nodes (15): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+7 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "KernelDev"
Cohesion: 0.16
Nodes (17): BinaryStat, binaryStatOf(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile(), BtrfsDevices(), btrfsSysfsDevices() (+9 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (15): recordTickCmd(), recordModel, clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits(), TestViewsFitReportedTerminalSize() (+7 more)

### Community 113 - "ConfinedEntry"
Cohesion: 0.12
Nodes (29): HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats(), newVerifyState() (+21 more)

### Community 115 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 116 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (60): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+52 more)

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

### Community 129 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.31
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "Candidate"
Cohesion: 0.13
Nodes (22): configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), cleanOrphanedFscrypt(), addManualDirectories(), groupCandidates(), pickDirectories(), pickFromCandidates() (+14 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "ListUsers"
Cohesion: 0.07
Nodes (29): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), newRecordFlow(), askDecrypt(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+21 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "MountVouch"
Cohesion: 0.13
Nodes (9): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), syncMap(), deleteKeys(), MountVouch, NewMountVouch(), NewTrustGuard() (+1 more)

### Community 143 - "BinaryEntry"
Cohesion: 0.10
Nodes (14): keptEntries, snapGuard, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), WithEagerPopulate(), diffBinaries() (+6 more)

### Community 144 - "DeployInstalledBinary"
Cohesion: 0.09
Nodes (30): installGUI(), buildBinaryIfNeeded(), installBinaryOnly(), installConfig(), installServices(), mustCwd(), revertSystemFiles(), applyUpdate() (+22 more)

### Community 145 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 148 - "New"
Cohesion: 0.08
Nodes (48): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+40 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.30
Nodes (13): appProtectorSet(), isAppProtector(), orphanedPolicies(), orphanedProtectors(), SelectOrphans(), appProtector(), foreignProtector(), policy() (+5 more)

### Community 150 - "netGuardModel"
Cohesion: 0.11
Nodes (11): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), protoString() (+3 more)

### Community 152 - "inodeChain"
Cohesion: 0.16
Nodes (19): fdStat, rootHandle, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent() (+11 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.06
Nodes (49): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), markReady(), newPinGeneration(), notifySystemdReady() (+41 more)

### Community 155 - "headless.go"
Cohesion: 0.23
Nodes (10): stagedBinary, hashFile(), hashReader(), stageRelease(), Check(), CheckGUI(), inspect(), buildTiny() (+2 more)

### Community 156 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 157 - "runBackend"
Cohesion: 0.12
Nodes (16): ParseEventsFlag(), reserveStdout(), runBackend(), LiveModeAvailable(), TestValidateFlagsEditConfig(), runEditProtected(), validateFlags(), TestValidateForwardFlags() (+8 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 163 - "Vault"
Cohesion: 0.10
Nodes (23): checkFreshDir(), openDirNoFollow(), renewDir(), FileModeToUnixMode(), Vault, IsBinaryContent(), Open(), inodeOf() (+15 more)

### Community 164 - "EventType"
Cohesion: 0.10
Nodes (13): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), DecodeBpfEvent() (+5 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.05
Nodes (47): actionBar, bleedLayout, cell, column, counterStrip, dataTable, eventFeed, feedAction (+39 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (39): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+31 more)

### Community 167 - "ValidateTargets"
Cohesion: 0.20
Nodes (8): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), Existing, ExistingRule

### Community 168 - "SanitizeText"
Cohesion: 0.12
Nodes (19): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), formatElapsed(), orNone(), printBinary(), printRecording() (+11 more)

### Community 170 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 176 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 177 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 178 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 179 - "finishBinary"
Cohesion: 0.14
Nodes (20): appletLine(), beyondEvents(), classifyGroup(), classifyLibs(), finishBinary(), innermostResource(), isShellHelper(), libCovered() (+12 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "secureResources"
Cohesion: 0.24
Nodes (8): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `--record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.11
Nodes (32): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+24 more)

### Community 204 - "runNetworkMonitor"
Cohesion: 0.18
Nodes (8): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NewNetModel(), NetworkMonitorUseCase

### Community 205 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "applyDiffAdditions"
Cohesion: 0.07
Nodes (43): confirmMasterKeyOverwrite(), runGenKey(), askBunTmpdirUsers(), gatherPerUserSetup(), addedResources(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), TestAddedResourcesSkipsAlreadyGuardedSSH() (+35 more)

### Community 210 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "diffModel"
Cohesion: 0.14
Nodes (3): diffModel, editorModel, changelogModel

### Community 214 - "Prereq"
Cohesion: 0.20
Nodes (6): classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 215 - "selection_test.go"
Cohesion: 0.25
Nodes (6): cellOf(), repeat(), selKey(), TestSelectionLargeDelete(), TestSelectionMouseDrag(), TestSelectionReadOnly()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 222 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 224 - "time.Time"
Cohesion: 0.12
Nodes (13): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+5 more)

### Community 225 - "main_test.go"
Cohesion: 0.11
Nodes (7): IntegrationSuite, IntegrationSuite, infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 227 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 228 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 236 - "fakeVault"
Cohesion: 0.36
Nodes (4): TestDaemonUseCaseGroupedEncryptionRootsDeduplicated(), TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 237 - "serve.go"
Cohesion: 0.08
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 239 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 240 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 241 - "Execute"
Cohesion: 0.25
Nodes (7): Requested(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand(), main()

### Community 243 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 244 - "record"
Cohesion: 0.29
Nodes (3): fixtureExe(), record(), TestObserveKeysByKernelExeAndSkipsSelf()

### Community 245 - "lockAndDeprovision"
Cohesion: 0.33
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 246 - "github.com/google/fscrypt/actions.Context"
Cohesion: 0.40
Nodes (4): modifiedContextWithSource(), CleanOrphans(), listPolicies(), listProtectors()

### Community 248 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 249 - "testRecording"
Cohesion: 0.47
Nodes (4): TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection(), testRecording(), TestRenderChecksTheSectionsMaterialize()

### Community 250 - "installHooks"
Cohesion: 0.50
Nodes (4): installHooks(), releaseBuild(), Run(), staged()

### Community 251 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 252 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

## Knowledge Gaps
- **32 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 616 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **57 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `RunWith`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09267296662254645 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018984547461368653 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `RunWith` to `daemon/daemon.go`, `edit-protected`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._