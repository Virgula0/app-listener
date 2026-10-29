# Graph Report - app-listener  (2026-10-29)

## Corpus Check
- 534 files · ~1,754,639 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6941 nodes · 30809 edges · 249 communities (208 shown, 41 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2865 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `524efbd1`
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
- tempGrantSetup
- New
- protectedPage
- GuardInodeKey
- daemonconfig.go
- root.go
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildTrustedSet
- guard.bpf.c
- TrustGuard
- changelog.go
- FileTarget
- golang.org/x/sys/unix.Stat_t
- github.com/charmbracelet/bubbletea.Cmd
- TestMonitorUseCaseLifecycle
- binaryList
- Guard
- picker
- orphans_test.go
- IntegrationSuite
- Message
- controlServer
- .posAt
- runDaemon
- IntegrationSuite
- absPath
- monitorUnitTest
- go_pkg_strings
- fileedit_model_test.go
- EventType
- stdio
- update_test.go
- fadeScroll
- Vault
- editView
- Client
- watchGroup
- SanitizeText
- .apply
- NewDumpHook
- mutedText
- engine
- guardModel
- audit.go
- Recorder
- inode_dev
- systemd.go
- UI
- unix_connect.c
- wm.go
- IntegrationSuite
- Config
- liveSession
- .showSized
- RefreshSection
- pickfs.go
- fcntl
- bpfstats/main.go
- globBuilder
- codecheck/main.go
- GuardEvent
- BinaryEntry
- process_vm_readv.c
- parseGuardEvent
- maintainer/plan.go
- LibraryClosure
- rcUser
- Monitor
- journal
- ListUsers
- removeFacultatives
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- CandidateDir
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pinstate.go
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- inspector_probe.c
- recordops.c
- os.FileInfo
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- update/update.go
- github.com/stretchr/testify/suite.Suite
- codeedit.go
- recordFlow
- IntegrationSuite
- makeReloadHandler
- newTestVetter
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- KernelDev
- WebSocket /ws client
- Selection
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- go_pkg_os
- updatePage
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- serve.go
- NetworkMonitor
- .populateWatchInodes
- Ledger
- runBackend
- expandPlaceholders
- .TestReconcileInodesEvictsStale
- DumpHook
- fakeStore
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- buildOneGuard
- runNetworkGuard
- planUpdaters
- gui/session.go
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- restore.go
- IntegrationSuite
- NetEventType
- networkPage
- shell
- cursor.go
- watchSet
- githubRelease
- deploy
- Limitations
- net_tester/main.go
- IntegrationSuite
- compileGlobNames
- safeio.go
- vaultfs_test.go
- RawTarget
- TestMain
- Compatibility
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- Binary
- dialog_test.go
- edit-protected
- github.com/spf13/cobra.Command
- kernel-prof.sh
- CopyTreeWithProgress
- pathCache
- .FilterExistingWhitelist
- time.Time
- Daemon configuration
- fileedit_symlink_test.go
- hashCache
- Desktop GUI
- IntegrationSuite
- logLine
- daemon: permanent protection with encryption at rest
- .eventTypeAllowed
- Active
- netinject.c
- procstats.go
- .awaitDaemonUpOrDead
- NetworkMonitorRepository
- VouchedDevs
- sync.Mutex
- unit.go
- IntegrationSuite
- openBunTmp
- changelog_test.go
- AddServeFlags
- Existing
- supersede_probe.c
- Prereq
- BinaryStat
- logGrid
- configedit_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- loadDaemonConfig
- FetchLatest
- vaultfs.go
- prctl_gate
- WithBottomBar
- Troubleshooting
- mountTable
- IntegrationSuite
- parsePasswd
- bleedLayout
- ptrace_seize.c
- Execute
- headless/headless.go
- monitor: see what touches a path
- downloadFile
- .exeKeys
- .runEditorHarness
- TestCleanRel

## God Nodes (most connected - your core abstractions)
1. `Guard` - 123 edges
2. `Config` - 118 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 97 edges
5. `absPath()` - 91 edges
6. `Load()` - 90 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 67 edges

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

## Communities (249 total, 41 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.17
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (163): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+155 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.13
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+10 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (27): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+19 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "tempGrantSetup"
Cohesion: 0.25
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 11 - "New"
Cohesion: 0.07
Nodes (51): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+43 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (24): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordView, trustPage (+16 more)

### Community 13 - "GuardInodeKey"
Cohesion: 0.03
Nodes (94): scriptIdentity(), checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), VettedInode, Multicall, MulticallKind (+86 more)

### Community 14 - "daemonconfig.go"
Cohesion: 0.36
Nodes (6): recordTick, rejectSymlinkedComponents(), validateWatchTarget(), isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 15 - "root.go"
Cohesion: 0.10
Nodes (12): releaseBuild(), Run(), TestEveryInstallModeMapped(), closeDumpFile(), postRun(), stageItem, stageState, textAnswer (+4 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "navItem"
Cohesion: 0.04
Nodes (15): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+7 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (10): Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (76): exe_refused(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+68 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (18): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust() (+10 more)

### Community 28 - "FileTarget"
Cohesion: 0.13
Nodes (9): guardPill(), FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps(), stopPumped() (+1 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (25): placedWalk, treeJudge, mountinfoDev(), SuperblockDev(), TestMountinfoDev(), confineBelow(), inUserHome(), vouchedFd() (+17 more)

### Community 30 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (16): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+8 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "binaryList"
Cohesion: 0.10
Nodes (13): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+5 more)

### Community 33 - "Guard"
Cohesion: 0.05
Nodes (19): binaryVerifyState, deferredBinary, exeRow, rootHandle, supersededBinary, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded() (+11 more)

### Community 34 - "picker"
Cohesion: 0.14
Nodes (6): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "orphans_test.go"
Cohesion: 0.18
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - ".posAt"
Cohesion: 0.10
Nodes (15): lineStyles, viewportState, fg(), TestDaemonConfTokens(), TestTokenizeKeepsLineCount(), tokenClassOf(), baseStyle(), gutterLabel() (+7 more)

### Community 40 - "runDaemon"
Cohesion: 0.08
Nodes (37): CheckBPFLSM(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), newPinGeneration() (+29 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (7): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "go_pkg_strings"
Cohesion: 0.06
Nodes (18): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin() (+10 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "EventType"
Cohesion: 0.07
Nodes (19): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+11 more)

### Community 48 - "update_test.go"
Cohesion: 0.10
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (18): appTheme, fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newTestFadeScroll(), settle() (+10 more)

### Community 50 - "Vault"
Cohesion: 0.22
Nodes (4): Vault, Open(), WriteFileInPlace(), WriteFileKeepMeta()

### Community 51 - "editView"
Cohesion: 0.12
Nodes (6): codeEntry, editView, Entry, fileSize(), joinRel(), relOf()

### Community 52 - "Client"
Cohesion: 0.09
Nodes (13): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+5 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 54 - "SanitizeText"
Cohesion: 0.07
Nodes (24): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), startTrustGuard() (+16 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 57 - "mutedText"
Cohesion: 0.05
Nodes (48): actionGate, binaryRow, binSource, filesPage, handButton, homePage, iconTile, installAction (+40 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (30): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin() (+22 more)

### Community 59 - "guardModel"
Cohesion: 0.06
Nodes (34): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), listenForGuardEvents(), tickGuardStats(), appendCapped(), clipToWidth() (+26 more)

### Community 60 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.12
Nodes (18): maskFor(), newBinary(), drain(), Recorder, newRecorder(), opNames(), selfExe(), sortedKeys() (+10 more)

### Community 62 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+15 more)

### Community 63 - "systemd.go"
Cohesion: 0.06
Nodes (63): offerDaemonStart(), runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+55 more)

### Community 64 - "UI"
Cohesion: 0.13
Nodes (15): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+7 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "Config"
Cohesion: 0.08
Nodes (50): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+42 more)

### Community 69 - "liveSession"
Cohesion: 0.17
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+4 more)

### Community 70 - ".showSized"
Cohesion: 0.11
Nodes (13): installUI, multiAnswer, sizedDialog, scrollable(), selectable(), newFadeScroll(), ask(), diffText() (+5 more)

### Community 71 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 72 - "pickfs.go"
Cohesion: 0.21
Nodes (15): listOpts, pickEntry, treeSearch, lstatOf(), outsideHardlinks(), TestSearchTreeFollowsNoSymlink(), plural(), entryIsDir() (+7 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "globBuilder"
Cohesion: 0.19
Nodes (10): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestGlobReservationsLibTrusted() (+2 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "GuardEvent"
Cohesion: 0.04
Nodes (26): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+18 more)

### Community 78 - "BinaryEntry"
Cohesion: 0.07
Nodes (23): keptEntries, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), canonicalPaths(), diffBinaries(), TestDiffBinaries() (+15 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 81 - "maintainer/plan.go"
Cohesion: 0.05
Nodes (52): SystemTrusted(), clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath() (+44 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 83 - "rcUser"
Cohesion: 0.12
Nodes (41): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+33 more)

### Community 84 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 85 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 86 - "ListUsers"
Cohesion: 0.09
Nodes (25): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), runOfflineEdit() (+17 more)

### Community 87 - "removeFacultatives"
Cohesion: 0.09
Nodes (32): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+24 more)

### Community 88 - "GenerateConf"
Cohesion: 0.05
Nodes (61): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH() (+53 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.10
Nodes (5): daemonPage, runHost, runView, captureLogs(), orDefault()

### Community 92 - "CandidateDir"
Cohesion: 0.11
Nodes (29): nextConfig(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+21 more)

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.16
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (9): lsmHook, PumpRingbuf(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook, NetGuard (+1 more)

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
Cohesion: 0.14
Nodes (27): HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats(), publishHash() (+19 more)

### Community 101 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.09
Nodes (12): snapNetGuard, snapNetMonitor, NetEvent, NetRecord, FormatAddr(), Ntohs(), reasonText(), NetGuardEvent (+4 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "os.FileInfo"
Cohesion: 0.17
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), TestCopyXattrsPreservesUserAttributes(), copyXattrs()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "update/update.go"
Cohesion: 0.15
Nodes (17): channelOf(), checkUpdatePreconditions(), deployGUI(), InstalledGUIVersion(), InstalledVersion(), Latest(), parseVerPart(), readInstalledVersion() (+9 more)

### Community 111 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.15
Nodes (5): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), harnessSuite, editorHarness

### Community 112 - "codeedit.go"
Cohesion: 0.12
Nodes (39): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+31 more)

### Community 113 - "recordFlow"
Cohesion: 0.16
Nodes (15): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+7 more)

### Community 115 - "makeReloadHandler"
Cohesion: 0.12
Nodes (14): newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), controlManager, WarnInodeOccupancy() (+6 more)

### Community 116 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "KernelDev"
Cohesion: 0.15
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "Selection"
Cohesion: 0.19
Nodes (7): pos, Highlighter, lexerFor(), shebangLexer(), cursorPos(), Selection, offset()

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

### Community 130 - "go_pkg_os"
Cohesion: 0.07
Nodes (54): confirmMasterKeyOverwrite(), runGenKey(), asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureMasterKey(), fileSHA256() (+46 more)

### Community 131 - "updatePage"
Cohesion: 0.20
Nodes (9): installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), newUpdatePage(), boolIndex() (+1 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.13
Nodes (16): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), grantNewBinaries() (+8 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.12
Nodes (11): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), startCatalogRefresh() (+3 more)

### Community 143 - "serve.go"
Cohesion: 0.07
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 144 - "NetworkMonitor"
Cohesion: 0.15
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (24): candidates(), confirmAll(), describe(), record(), applySchema(), Ledger, Pending, Source (+16 more)

### Community 147 - "runBackend"
Cohesion: 0.10
Nodes (21): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), pingInterval() (+13 more)

### Community 148 - "expandPlaceholders"
Cohesion: 0.14
Nodes (7): expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 149 - ".TestReconcileInodesEvictsStale"
Cohesion: 0.22
Nodes (10): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+2 more)

### Community 150 - "DumpHook"
Cohesion: 0.24
Nodes (5): logHook, formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 152 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (28): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), newFindInput() (+20 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (44): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard() (+36 more)

### Community 155 - "runNetworkGuard"
Cohesion: 0.07
Nodes (34): DrainUntilSignal(), CheckEBPF(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless() (+26 more)

### Community 156 - "planUpdaters"
Cohesion: 0.08
Nodes (19): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+11 more)

### Community 157 - "gui/session.go"
Cohesion: 0.18
Nodes (16): invokingIDs(), ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+8 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "restore.go"
Cohesion: 0.21
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 164 - "NetEventType"
Cohesion: 0.15
Nodes (12): NetEventType, eventsOrAll(), NetEventTypes(), ParseNetEventType(), WithSuspectsFrom(), eventsetSummary(), eventTypeKey(), modeLabel() (+4 more)

### Community 165 - "networkPage"
Cohesion: 0.07
Nodes (39): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+31 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (66): anyDaemonRunning(), fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapEngines (+58 more)

### Community 167 - "cursor.go"
Cohesion: 0.07
Nodes (27): buttonCursor, checkGrid, handCheck, handCheckGroup, handRadioGroup, handSelect, recordReview, handCursor() (+19 more)

### Community 168 - "watchSet"
Cohesion: 0.13
Nodes (15): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+7 more)

### Community 169 - "githubRelease"
Cohesion: 0.27
Nodes (16): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+8 more)

### Community 170 - "deploy"
Cohesion: 0.18
Nodes (18): applyUpdate(), archAssets(), assetsFor(), confirmUnit(), deploy(), downloadAndVerify(), downloadReleaseFiles(), fetchGUI() (+10 more)

### Community 172 - "Limitations"
Cohesion: 0.20
Nodes (10): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+2 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 176 - "safeio.go"
Cohesion: 0.23
Nodes (14): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), AtomicWriteAt(), CreateExclNoFollow() (+6 more)

### Community 177 - "vaultfs_test.go"
Cohesion: 0.20
Nodes (14): IsBinaryContent(), inodeOf(), mustRead(), mustWrite(), openRootT(), TestFileModeToUnixModeRoundTrip(), TestIsBinaryContent(), TestOpenVaultRefusesSymlink() (+6 more)

### Community 178 - "RawTarget"
Cohesion: 0.22
Nodes (12): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+4 more)

### Community 179 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

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

### Community 194 - "Binary"
Cohesion: 0.14
Nodes (18): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+10 more)

### Community 195 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "github.com/spf13/cobra.Command"
Cohesion: 0.25
Nodes (13): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink() (+5 more)

### Community 199 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 202 - "time.Time"
Cohesion: 0.08
Nodes (35): LockedError, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+27 more)

### Community 203 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 204 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 205 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "logLine"
Cohesion: 0.25
Nodes (5): logLine, journalLine(), TestJournalLine(), logRow(), logStyle()

### Community 210 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.15
Nodes (13): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+5 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "Active"
Cohesion: 0.20
Nodes (10): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), Active() (+2 more)

### Community 213 - "netinject.c"
Cohesion: 0.32
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 214 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 216 - "NetworkMonitorRepository"
Cohesion: 0.18
Nodes (4): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 218 - "sync.Mutex"
Cohesion: 0.13
Nodes (9): refusedReplacements, fyneLogFilter, TestRefusedReplacements_ReportsEachInodeOnce(), newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity (+1 more)

### Community 219 - "unit.go"
Cohesion: 0.33
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 221 - "openBunTmp"
Cohesion: 0.53
Nodes (6): fileID(), openBunTmp(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 222 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 223 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 224 - "Existing"
Cohesion: 0.18
Nodes (8): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), Existing, ExistingRule

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 227 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 229 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 231 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 232 - "FetchLatest"
Cohesion: 0.33
Nodes (9): installHooks(), staged(), FetchLatest(), fetchPair(), fetchRelease(), FetchVerified(), Release, TestFetchReleaseByTag() (+1 more)

### Community 233 - "vaultfs.go"
Cohesion: 0.22
Nodes (8): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsOctalMode(), ParseMode(), TestDefaultNewFileMode(), TestIsOctalMode(), Unopened()

### Community 234 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 235 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.39
Nodes (5): MountInfo, MountByID(), ParseMountinfo(), ReadMountinfo(), mountTable

### Community 239 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 241 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 242 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 243 - "headless/headless.go"
Cohesion: 0.43
Nodes (5): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 245 - "downloadFile"
Cohesion: 0.33
Nodes (4): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 640 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.044) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018644650842793258 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `TrustGuard`, `engine`, `time.Duration`, `ConfinedEntry`, `runDaemon`, `guardUnitTest`, `time.Time`, `sync.Mutex`, `GuardInodeKey`, `GuardEvent`, `BinaryEntry`, `makeReloadHandler`, `.TestReconcileInodesEvictsStale`, `buildOneGuard`, `planUpdaters`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09124209124209125 - nodes in this community are weakly interconnected._
- **Why does `backend` connect `Message` to `go_pkg_os`, `time.Duration`, `Vault`, `liveSession`, `recordFlow`, `Vault`, `runBackend`, `Client`, `TestCleanRel`, `sync.Mutex`, `planUpdaters`, `backend_test.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Vault` be split into smaller, more focused modules?**
  _Cohesion score 0.13446969696969696 - nodes in this community are weakly interconnected._