# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 486 files · ~1,715,324 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6490 nodes · 27910 edges · 216 communities (169 shown, 47 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2600 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1bdb01f6`
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
- recorder_test.go
- policyEditor
- protectedPage
- codeedit.go
- CandidateDir
- editView
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- bpfstats/main.go
- guard.bpf.c
- TrustGuard
- github.com/charmbracelet/bubbles/textarea.Model
- github.com/charmbracelet/bubbletea.Cmd
- shellrc.go
- Config
- runMonitor
- ui.go
- NetEventType
- picker
- fileedit_model_test.go
- time.Duration
- Message
- controlServer
- fakeGuardRepo
- ConfinedEntry
- IntegrationSuite
- absPath
- monitorUnitTest
- wizard.go
- loadDaemonConfig
- selection.go
- stdio
- netModel
- fadeScroll
- newTestTarget
- cursor.go
- proto.go
- watchGroup
- navItem
- Plan
- RawTarget
- mutedText
- GuardInodeKey
- os.File
- User
- Target
- exe_supersede.h
- BinaryEntry
- RefreshSection
- pinstate.go
- wm.go
- IntegrationSuite
- Binary
- hashFileShared
- fcntl
- model
- Guard
- sync.Mutex
- Candidate
- catalogRefresher
- layout.go
- auditGroupCoverage
- process_vm_readv.c
- time.Time
- changelog_test.go
- mayUpdate
- update/update.go
- net_tester/main.go
- appTheme
- render.go
- recordFlow
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- inodeChain
- Monitor
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- io.Writer
- changelog.go
- protectedShell
- Well-known bypass classes
- installPage
- jit_provenance.c
- inspector_probe.c
- recordops.c
- fakeDaemon
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- makeReloadHandler
- mountTable
- netGuardModel
- parseGuardEvent
- IntegrationSuite
- logGrid
- vaultfs.go
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
- install/install.go
- trace-app-libs.sh
- IntegrationSuite
- dialLiveSession
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- GuardEvent
- IntegrationSuite
- Choice
- newTestVetter
- DumpHook
- diffModel
- fscrypt.go
- fakeStore
- New
- IntegrationSuite
- SanitizeText
- IntegrationSuite
- Installation
- backend_test.go
- graphify-refresh.sh
- hashCache
- EventType
- feed.go
- shell
- github.com/stretchr/testify/suite.Suite
- Ledger
- taint_victim.c
- find_test.go
- Limitations
- Troubleshooting
- IntegrationSuite
- copyRegularAt
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
- Execute
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- NewDumpHook
- pathCache
- NewEventFanout
- DaemonEvent
- Desktop GUI
- IntegrationSuite
- .eventTypeAllowed
- bleedLayout
- TestMain
- LibraryClosure
- runNetworkGuard
- Daemon configuration
- watchSet
- CleanupStalePins
- .runEditorHarness
- inode_dev.h
- desktop.go
- Serve

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `GuardInodeKey` - 99 edges
4. `Load()` - 89 edges
5. `shell` - 86 edges
6. `absPath()` - 83 edges
7. `IntegrationSuite` - 82 edges
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

## Communities (216 total, 47 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (23): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+15 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (144): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload() (+136 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.07
Nodes (26): encryptDirectories(), secureResources(), verifyEncryptionState(), askEncryption(), deprovisionKind, harnessSuite, migrationTarget, isFileVaultCiphertext() (+18 more)

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
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "recorder_test.go"
Cohesion: 0.34
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 11 - "policyEditor"
Cohesion: 0.13
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (24): backendLink, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage, recordForm (+16 more)

### Community 13 - "codeedit.go"
Cohesion: 0.15
Nodes (27): class, editorModel, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+19 more)

### Community 14 - "CandidateDir"
Cohesion: 0.05
Nodes (48): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+40 more)

### Community 15 - "editView"
Cohesion: 0.12
Nodes (6): codeEntry, editView, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (11): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), chownUser, dirEntry, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (55): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+47 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (13): ownerPaths(), Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+5 more)

### Community 23 - "auth_test.go"
Cohesion: 0.14
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (27): GlobChild, statFunc, tempGrant, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), singleBinaryPaths() (+19 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.18
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (14): recordTickCmd(), recordModel, NewDaemonModel(), clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits() (+6 more)

### Community 29 - "shellrc.go"
Cohesion: 0.07
Nodes (58): setupBunTmpdir(), refreshCatalogConfig(), installFile(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt() (+50 more)

### Community 30 - "Config"
Cohesion: 0.06
Nodes (46): newBinaryVetter(), openBinaryVetter(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure() (+38 more)

### Community 31 - "runMonitor"
Cohesion: 0.16
Nodes (9): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), MonitorRepository, NewModel(), MonitorUseCase, NewMonitorUseCase() (+1 more)

### Community 32 - "ui.go"
Cohesion: 0.09
Nodes (23): nextConfig(), TestConfirmUpdateRoutesToActiveUI(), InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks (+15 more)

### Community 33 - "NetEventType"
Cohesion: 0.05
Nodes (30): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+22 more)

### Community 34 - "picker"
Cohesion: 0.06
Nodes (45): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+37 more)

### Community 35 - "fileedit_model_test.go"
Cohesion: 0.19
Nodes (26): newFileEditModel(), RunFileEditor(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow() (+18 more)

### Community 36 - "time.Duration"
Cohesion: 0.12
Nodes (11): IntegrationSuite, pooledContainer, formatAgo(), formatClock(), New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod() (+3 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), backend, backendOp, backendState, fakeLink (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), TestRunConfigEditOverTheSocket(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+15 more)

### Community 39 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (24): SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock(), ResolveTempBinary() (+16 more)

### Community 40 - "ConfinedEntry"
Cohesion: 0.15
Nodes (29): markReady(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), ConfinedEntry(), HashStats(), hashedReused() (+21 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.15
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "wizard.go"
Cohesion: 0.08
Nodes (34): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+26 more)

### Community 45 - "loadDaemonConfig"
Cohesion: 0.20
Nodes (10): loadDaemonConfig(), noteReloadedSettings(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes() (+2 more)

### Community 46 - "selection.go"
Cohesion: 0.20
Nodes (5): pos, reposition(), cursorPos(), Selection, offset()

### Community 48 - "netModel"
Cohesion: 0.17
Nodes (11): runNetworkMonitor(), runTUI(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel(), tickNetStats() (+3 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 51 - "cursor.go"
Cohesion: 0.10
Nodes (14): binaryList, buttonCursor, newBinaryList(), handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm() (+6 more)

### Community 52 - "proto.go"
Cohesion: 0.06
Nodes (24): BinaryCheck, Client, Entry, Event, File, Progress, PwState, Reader (+16 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 54 - "navItem"
Cohesion: 0.03
Nodes (19): clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar, tabButton (+11 more)

### Community 55 - "Plan"
Cohesion: 0.09
Nodes (21): BinaryRule, InsertSectionsBeforeLibraries(), dedupe(), Plan, libOfChosen(), materialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale() (+13 more)

### Community 56 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (42): binaryRow, handButton, homePage, iconTile, latestRelease, panel, pill, statusTile (+34 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (52): engine, refusedReplacements, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+44 more)

### Community 59 - "os.File"
Cohesion: 0.03
Nodes (77): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp (+69 more)

### Community 60 - "User"
Cohesion: 0.06
Nodes (50): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), addedResources() (+42 more)

### Community 61 - "Target"
Cohesion: 0.11
Nodes (28): appletLine(), classifyGroup(), classifyLibs(), innermostResource(), libCovered(), newBinary(), verifyExe(), worldWritableDir() (+20 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "BinaryEntry"
Cohesion: 0.10
Nodes (14): keptEntries, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), diffBinaries(), TestDiffBinaries() (+6 more)

### Community 64 - "RefreshSection"
Cohesion: 0.10
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 65 - "pinstate.go"
Cohesion: 0.11
Nodes (32): confirmMasterKeyOverwrite(), resolveConfigPath(), runGenKey(), runLockdown(), runOneShotMode(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+24 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "Binary"
Cohesion: 0.09
Nodes (26): formatElapsed(), lines(), orNone(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording() (+18 more)

### Community 69 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (28): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, rootHandle, tempMaskOp (+20 more)

### Community 73 - "sync.Mutex"
Cohesion: 0.07
Nodes (22): filesPage, Target, guardPill(), newActivity(), newExeResolver(), Files, FileSpec, FileTarget (+14 more)

### Community 74 - "Candidate"
Cohesion: 0.25
Nodes (14): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+6 more)

### Community 76 - "catalogRefresher"
Cohesion: 0.15
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), startCatalogRefresh(), catalogRefresher (+1 more)

### Community 77 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 78 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+4 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "time.Time"
Cohesion: 0.27
Nodes (5): LockedError, sweepID, sweepMemo, Progress, Seen

### Community 81 - "changelog_test.go"
Cohesion: 0.24
Nodes (7): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 82 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 83 - "update/update.go"
Cohesion: 0.04
Nodes (91): installHooks(), staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), tags() (+83 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+6 more)

### Community 87 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 88 - "liveSession"
Cohesion: 0.22
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 92 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.05
Nodes (34): fdStat, snapNetGuard, physicalParent(), statFD(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices() (+26 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "io.Writer"
Cohesion: 0.31
Nodes (8): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, candidate

### Community 100 - "changelog.go"
Cohesion: 0.15
Nodes (16): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 101 - "protectedShell"
Cohesion: 0.27
Nodes (8): checkGrid, newCheckGrid(), isOwnExecutable(), newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestRecordReviewExtendsAnExistingSection()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "installPage"
Cohesion: 0.11
Nodes (13): actionGate, binSource, installAction, installedBinary, installPage, stagedBinary, hashFile(), hashReader() (+5 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "makeReloadHandler"
Cohesion: 0.29
Nodes (4): newControlManager(), startControlManager(), makeReloadHandler(), controlManager

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "netGuardModel"
Cohesion: 0.10
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 113 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 116 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (59): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+51 more)

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
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - ".launch"
Cohesion: 0.31
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "install/install.go"
Cohesion: 0.05
Nodes (77): applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), deploy(), ensureInstalledBinary() (+69 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "dialLiveSession"
Cohesion: 0.05
Nodes (42): ParseEventsFlag(), auditAfterEdit(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), TestValidateFlagsEditConfig(), applyEdit() (+34 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "GuardEvent"
Cohesion: 0.08
Nodes (18): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), GuardEvent (+10 more)

### Community 144 - "Choice"
Cohesion: 0.50
Nodes (3): MultiSelectKeymap(), Choice, MultiSelect()

### Community 145 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 147 - "DumpHook"
Cohesion: 0.16
Nodes (8): logHook, logLine, logRow(), logStyle(), formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 149 - "fscrypt.go"
Cohesion: 0.05
Nodes (52): ensureMasterKey(), confirmRunPrereq(), resolveFilesystemPrereqs(), invalidPassword, writeTempWithMetadata(), checkKeyLen(), classifySetupError(), classifySupportError() (+44 more)

### Community 152 - "New"
Cohesion: 0.09
Nodes (45): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+37 more)

### Community 154 - "SanitizeText"
Cohesion: 0.04
Nodes (76): CheckEBPF(), multicallAdmissible(), backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), configureDaemonLogging(), eventFilterOptions() (+68 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 162 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 164 - "EventType"
Cohesion: 0.07
Nodes (19): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+11 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (43): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (48): pageDef, pageID, procStats, reel, resizeWatch, shell, sizedDialog, snapEngines (+40 more)

### Community 168 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 170 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 177 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.14
Nodes (19): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+11 more)

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

### Community 195 - "Execute"
Cohesion: 0.29
Nodes (7): Requested(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand(), main()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `--record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (36): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+28 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 215 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 219 - "runNetworkGuard"
Cohesion: 0.05
Nodes (36): CheckBPFLSM(), ParseNetEventsFlag(), runBPFCheck(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard() (+28 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 224 - "watchSet"
Cohesion: 0.13
Nodes (15): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+7 more)

### Community 227 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 229 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 230 - "desktop.go"
Cohesion: 0.07
Nodes (32): releaseBuild(), Run(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+24 more)

### Community 237 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

## Knowledge Gaps
- **31 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 600 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11680465428507496 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `wizard.go`, `Message`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02059453918553723 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `wizard.go` to `recorder_test.go`, `edit-protected`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._